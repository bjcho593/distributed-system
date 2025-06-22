using Microsoft.AspNetCore.Mvc;

namespace login_service.Controllers
{
    [ApiController]
    [Route("login")]
    public class LoginController : ControllerBase
    {
        [HttpPost("auth")]
        public IActionResult Login([FromBody] LoginRequest request)
        {
            // Aquí podrías conectar a MySQL para validar usuario
            if (request.Username == "admin" && request.Password == "admin123")
                return Ok(new { message = "Login exitoso" });

            return Unauthorized(new { message = "Credenciales inválidas" });
        }
    }

    public class LoginRequest
    {
        public string Username { get; set; }
        public string Password { get; set; }
    }
}
