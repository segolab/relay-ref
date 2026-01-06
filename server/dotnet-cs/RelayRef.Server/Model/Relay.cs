namespace RelayRef.Server.Model;

public sealed class Relay
{
    public Guid Id { get; init; }
    public string EventType { get; init; } = "";
    public Destination Destination { get; init; } = new();
    public object Payload { get; init; } = default!;
    public Dictionary<string, string>? Metadata { get; init; }
    public string Status { get; init; } = "queued";
    public DateTime CreatedAt { get; init; }
    public DateTime? DeliveredAt { get; init; }
    public string? FailureReason { get; init; }
}
