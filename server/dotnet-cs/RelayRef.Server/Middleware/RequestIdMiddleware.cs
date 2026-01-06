namespace RelayRef.Server.Middleware;

public sealed class RequestIdMiddleware
{
    public const string HeaderName = "X-Request-ID";

    private readonly RequestDelegate _next;

    public RequestIdMiddleware(RequestDelegate next)
    {
        _next = next;
    }

    public async Task Invoke(HttpContext context)
    {
        var id = context.Request.Headers[HeaderName].FirstOrDefault()
                 ?? Guid.NewGuid().ToString();

        context.Items[HeaderName] = id;
        context.Response.Headers[HeaderName] = id;

        await _next(context);
    }

    public static string? Get(HttpContext ctx)
        => ctx.Items.TryGetValue(HeaderName, out var v) ? v as string : null;
}
