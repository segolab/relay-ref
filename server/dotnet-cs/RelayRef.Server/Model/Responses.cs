namespace RelayRef.Server.Model;

public sealed class ErrorResponse
{
    public string Code { get; init; } = "";
    public string Message { get; init; } = "";
    public object? Details { get; init; }
    public string? RequestId { get; init; }
}
