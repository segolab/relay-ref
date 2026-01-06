using RelayRef.Server.Configuration;
using RelayRef.Server.Middleware;
using RelayRef.Server.RateLimiting;
using RelayRef.Server.Store;

var builder = WebApplication.CreateBuilder(args);

var rawOptions = RelayEnvLoader.Load(builder.Configuration);

builder.Services.AddSingleton(rawOptions);

builder.WebHost.UseUrls(
    HttpAddrNormalizer.ToUrl(rawOptions.HttpAddr)
);

builder.Services.Configure<RelayOptions>(
    builder.Configuration.GetSection("Relay"));

builder.Services.AddSingleton<InMemoryRelayStore>();
builder.Services.AddSingleton<InMemoryIdempotencyStore>();
builder.Services.AddSingleton<TokenBucketLimiter>();

builder.Services.AddControllers();

var app = builder.Build();

// Middleware order = Diagram 02
app.UseMiddleware<ExceptionMiddleware>();
app.UseMiddleware<RequestIdMiddleware>();
app.UseMiddleware<ApiKeyMiddleware>();
app.UseMiddleware<RateLimitMiddleware>();

app.MapControllers();

app.MapGet("/healthz", () => Results.Ok());
app.MapGet("/readyz", () => Results.Ok());

app.Run();

public partial class Program { }
