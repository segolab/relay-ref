using Microsoft.Extensions.Configuration;

namespace RelayRef.Server.Configuration;

public static class RelayEnvLoader
{
    public static RelayOptions Load(IConfiguration config)
    {
        return new RelayOptions
        {
            HttpAddr = Get(config, "RELAY_HTTP_ADDR", ":8429"),
            ApiKeys = Get(config, "RELAY_API_KEYS", "dev-key")
                .Split(',', StringSplitOptions.RemoveEmptyEntries | StringSplitOptions.TrimEntries),

            MaxBodyBytes = GetInt(config, "RELAY_MAX_BODY_BYTES", 32_768),
            IdempotencyTtlSeconds = GetInt(config, "RELAY_IDEMPOTENCY_TTL_SECONDS", 3600),

            RateLimit = new RateLimitOptions
            {
                PostRps = GetDouble(config, "RELAY_LIMIT_POST_RPS", 10),
                PostBurst = GetInt(config, "RELAY_LIMIT_POST_BURST", 20),
                GetRps = GetDouble(config, "RELAY_LIMIT_GET_RPS", 50),
                GetBurst = GetInt(config, "RELAY_LIMIT_GET_BURST", 100),
            }
        };
    }

    private static string Get(IConfiguration c, string key, string def)
        => string.IsNullOrEmpty(c[key]) ? def : c[key]!;

    private static int GetInt(IConfiguration c, string key, int def)
        => int.TryParse(c[key], out var v) ? v : def;

    private static double GetDouble(IConfiguration c, string key, double def)
        => double.TryParse(c[key], out var v) ? v : def;
}

public static class HttpAddrNormalizer
{
    public static string ToUrl(string addr)
    {
        // ":8429" → "http://0.0.0.0:8429"
        if (addr.StartsWith(":"))
            return "http://0.0.0.0" + addr;

        // "0.0.0.0:8429" → "http://0.0.0.0:8429"
        if (!addr.StartsWith("http://") && !addr.StartsWith("https://"))
            return "http://" + addr;

        return addr;
    }
}
