use redis::AsyncCommands;
use std::env;
use dotenv::dotenv;
use tracing::{info, error};
use tracing_subscriber;

#[tokio::main]
async fn main() -> redis::RedisResult<()> {
    dotenv().ok();
    tracing_subscriber::fmt::init();

    let redis_url = env::var("REDIS_URL").expect("REDIS_URL not set in .env");
    let client = redis::Client::open(redis_url)?;
    let mut con = client.get_async_connection().await?;

    // Guardar un valor
    con.set("hello", "world").await?;
    // Recuperar el valor
    let val: String = con.get("hello").await?;
    info!("Valor desde Redis: {}", val);

    Ok(())
}
