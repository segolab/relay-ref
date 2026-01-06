using System.Security.Cryptography;
using System.Text;
using System.Text.Json;
using Microsoft.AspNetCore.Mvc;
using Microsoft.Extensions.Options;
using RelayRef.Server.Configuration;
using RelayRef.Server.Middleware;
using RelayRef.Server.Model;
using RelayRef.Server.Store;

namespace RelayRef.Server.Api;

[ApiController]
[Route("v1/relays")]
public sealed class RelayController : ControllerBase
{
    private readonly InMemoryRelayStore _store;
    private readonly InMemoryIdempotencyStore _idem;
    private readonly RelayOptions _options;

    public RelayController(
        InMemoryRelayStore store,
        InMemoryIdempotencyStore idem,
        IOptions<RelayOptions> options)
    {
        _store = store;
        _idem = idem;
        _options = options.Value;
    }

    [HttpPost]
    public IActionResult Create(
        [FromBody] CreateRelayRequest request,
        [FromHeader(Name = "Idempotency-Key")] string? idemKey)
    {
        var apiKey = ApiKeyMiddleware.ApiKey(HttpContext);

        var raw = JsonSerializer.Serialize(request);
        var hash = Convert.ToHexString(
            SHA256.HashData(Encoding.UTF8.GetBytes(raw)));

        Func<Relay> create = () =>
        {
            var relay = new Relay
            {
                Id = Guid.NewGuid(),
                EventType = request.EventType,
                Destination = request.Destination,
                Payload = request.Payload,
                Metadata = request.Metadata,
                CreatedAt = DateTime.UtcNow
            };
            _store.Create(relay);
            return relay;
        };

        Relay result;

        try
        {
            result = idemKey is null
                ? create()
                : _idem.GetOrCreate(
                    apiKey,
                    idemKey,
                    hash,
                    create,
                    TimeSpan.FromSeconds(_options.IdempotencyTtlSeconds));
        }
        catch (InvalidOperationException)
        {
            return Conflict(new ErrorResponse
            {
                Code = "idempotency_conflict",
                Message = "Idempotency key reused with different payload",
                RequestId = RequestIdMiddleware.Get(HttpContext)
            });
        }

        return Created($"/v1/relays/{result.Id}", result);
    }
}
