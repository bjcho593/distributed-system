use axum::{
    routing::get,
    extract::Query,
    response::Json,
    Router,
};
use std::{collections::HashMap, env, net::SocketAddr};
use dotenv::dotenv;
use redis::AsyncCommands;
use tracing::{info, error};
use tracing_subscriber;

#[tokio::main]
async fn main() {
    dotenv().ok();
    tracing_subscriber::fmt::init();

    let app = Router::new().route("/search", get(search_handler));

    let addr = SocketAddr::from(([0, 0, 0, 0], 3000));
    info!("Servidor escuchando en http://{}", addr);

    axum::Server::bind(&addr)
        .serve(app.into_make_service())
        .await
        .unwrap();
}

async fn search_handler(Query(params): Query<HashMap<String, String>>) -> Json<HashMap<String, String>> {
    let query = params.get("q").cloned().unwrap_or_default();
    let redis_url = env::var("REDIS_URL").unwrap();
    let client = redis::Client::open(redis_url).unwrap();
    let mut con = client.get_async_connection().await.unwrap();

    // Buscar la palabra en Redis
    let result: Option<String> = con.get(&query).await.unwrap_or(None);

    let mut response = HashMap::new();
    response.insert("query".to_string(), query.clone());
    response.insert("result".to_string(), result.unwrap_or("No encontrado".to_string()));
    
    Json(response)
}
