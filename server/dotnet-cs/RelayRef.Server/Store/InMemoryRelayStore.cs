using RelayRef.Server.Model;

namespace RelayRef.Server.Store;

public sealed class InMemoryRelayStore
{
    private readonly Dictionary<Guid, Relay> _byId = new();
    private readonly List<Guid> _order = new();
    private readonly ReaderWriterLockSlim _lock = new();

    public void Create(Relay relay)
    {
        _lock.EnterWriteLock();
        try
        {
            _byId[relay.Id] = relay;
            _order.Add(relay.Id);
        }
        finally
        {
            _lock.ExitWriteLock();
        }
    }

    public Relay? Get(Guid id)
    {
        _lock.EnterReadLock();
        try
        {
            return _byId.TryGetValue(id, out var r) ? r : null;
        }
        finally
        {
            _lock.ExitReadLock();
        }
    }
}
