namespace RelayRef.Server.Configuration;

public sealed class RelayOptions
{
    public string HttpAddr { get; init; } = ":8429";

    public string[] ApiKeys { get; init; } = ["dev-key"];

    public int MaxBodyBytes { get; init; } = 32_768;

    public int IdempotencyTtlSeconds { get; init; } = 3600;

    public RateLimitOptions RateLimit { get; init; } = new();
}

public sealed class RateLimitOptions
{
    public double PostRps { get; init; } = 10;
    public int PostBurst { get; init; } = 20;

    public double GetRps { get; init; } = 50;
    public int GetBurst { get; init; } = 100;
}
