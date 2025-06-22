var builder = WebApplication.CreateBuilder(args);
builder.Services.AddControllers();
var app = builder.Build();

app.MapGet("/login/status", () => Results.Ok("Login Service is running ✅"));
app.MapControllers();

app.Run();
