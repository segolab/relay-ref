using RelayRef.Server.Model;

namespace RelayRef.Server.Store;

public sealed class InMemoryIdempotencyStore
{
    private sealed record Entry(string Hash, Relay Relay, DateTime ExpiresAt);

    private readonly Dictionary<string, Entry> _entries = new();
    private readonly object _lock = new();

    public Relay GetOrCreate(
        string apiKey,
        string idemKey,
        string payloadHash,
        Func<Relay> factory,
        TimeSpan ttl)
    {
        var key = $"{apiKey}:{idemKey}";
        var now = DateTime.UtcNow;

        lock (_lock)
        {
            if (_entries.TryGetValue(key, out var e))
            {
                if (now > e.ExpiresAt)
                {
                    _entries.Remove(key);
                }
                else if (e.Hash != payloadHash)
                {
                    throw new InvalidOperationException("idempotency_conflict");
                }
                else
                {
                    return e.Relay;
                }
            }

            var relay = factory();
            _entries[key] = new Entry(payloadHash, relay, now.Add(ttl));
            return relay;
        }
    }
}
