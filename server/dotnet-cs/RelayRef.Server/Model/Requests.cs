using System.Text.Json;

namespace RelayRef.Server.Model;

public sealed class CreateRelayRequest
{
    public string EventType { get; init; } = "";
    public Destination Destination { get; init; } = new();
    public JsonElement Payload { get; init; }
    public Dictionary<string, string>? Metadata { get; init; }
}

public sealed class Destination
{
    public string Type { get; init; } = "";
    public string Url { get; init; } = "";
}
