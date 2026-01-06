using Microsoft.Extensions.Options;
using RelayRef.Server.Configuration;

namespace RelayRef.Server.Middleware;

public sealed class ApiKeyMiddleware
{
    private readonly RequestDelegate _next;
    private readonly HashSet<string> _allowed;

    public ApiKeyMiddleware(RequestDelegate next, IOptions<RelayOptions> options)
    {
        _next = next;
        _allowed = options.Value.ApiKeys.ToHashSet();
    }

    public async Task Invoke(HttpContext context)
    {
        if (!context.Request.Headers.TryGetValue("X-API-Key", out var key) ||
            !_allowed.Contains(key!))
        {
            context.Response.StatusCode = 401;
            return;
        }

        context.Items["apiKey"] = key!.ToString();
        await _next(context);
    }

    public static string ApiKey(HttpContext ctx)
        => (string)ctx.Items["apiKey"]!;
}
