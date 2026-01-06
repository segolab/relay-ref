using Microsoft.Extensions.Options;
using RelayRef.Server.Configuration;
using RelayRef.Server.RateLimiting;

namespace RelayRef.Server.Middleware;

public sealed class RateLimitMiddleware
{
    private readonly RequestDelegate _next;
    private readonly TokenBucketLimiter _limiter;

    public RateLimitMiddleware(RequestDelegate next, TokenBucketLimiter limiter)
    {
        _next = next;
        _limiter = limiter;
    }

    public async Task Invoke(HttpContext context)
    {
        var apiKey = ApiKeyMiddleware.ApiKey(context);
        var path = context.Request.Path.Value ?? "";
        var routeGroup = path.StartsWith("/v1/relays") && context.Request.Method == "POST"
            ? "post_relays"
            : "get_relays";

        var result = _limiter.Allow(apiKey, routeGroup);

        context.Response.Headers["RateLimit-Limit"] = result.Limit.ToString();
        context.Response.Headers["RateLimit-Remaining"] = result.Remaining.ToString();
        context.Response.Headers["RateLimit-Reset"] = result.ResetInSeconds.ToString();

        if (!result.Allowed)
        {
            context.Response.Headers["Retry-After"] = result.RetryAfterSeconds.ToString();
            context.Response.StatusCode = 429;
            return;
        }

        await _next(context);
    }
}
