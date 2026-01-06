namespace RelayRef.Server.RateLimiting;

public sealed class TokenBucketLimiter
{
    private sealed class Bucket
    {
        private double _tokens;
        private DateTime _last;
        private readonly double _rps;
        private readonly double _burst;
        private readonly object _lock = new();

        public Bucket(double rps, int burst)
        {
            _rps = rps;
            _burst = burst;
            _tokens = burst;
            _last = DateTime.UtcNow;
        }

        public RateLimitResult TryConsume()
        {
            lock (_lock)
            {
                var now = DateTime.UtcNow;
                var delta = (now - _last).TotalSeconds;
                _tokens = Math.Min(_burst, _tokens + delta * _rps);
                _last = now;

                if (_tokens >= 1)
                {
                    _tokens -= 1;
                    return RateLimitResult.Allow((int)_burst, (int)_tokens);
                }

                var retry = (int)Math.Ceiling((1 - _tokens) / _rps);
                return RateLimitResult.Deny((int)_burst, (int)_tokens, retry);
            }
        }
    }

    private readonly Dictionary<string, Bucket> _buckets = new();
    private readonly object _lock = new();

    public RateLimitResult Allow(string apiKey, string routeGroup)
    {
        lock (_lock)
        {
            if (!_buckets.TryGetValue(apiKey, out var bucket))
            {
                bucket = routeGroup == "post_relays"
                    ? new Bucket(10, 20)
                    : new Bucket(50, 100);
                _buckets[apiKey] = bucket;
            }

            return bucket.TryConsume();
        }
    }
}

public sealed record RateLimitResult(
    bool Allowed,
    int Limit,
    int Remaining,
    int ResetInSeconds,
    int RetryAfterSeconds)
{
    public static RateLimitResult Allow(int limit, int remaining)
        => new(true, limit, remaining, 0, 0);

    public static RateLimitResult Deny(int limit, int remaining, int retry)
        => new(false, limit, remaining, retry, retry);
}
