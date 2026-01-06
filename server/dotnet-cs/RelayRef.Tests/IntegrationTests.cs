using System.Net;
using System.Net.Http.Json;
using Microsoft.AspNetCore.Mvc.Testing;
using Xunit;

public class IntegrationTests : IClassFixture<WebApplicationFactory<Program>>
{
    private readonly HttpClient _client;

    public IntegrationTests(WebApplicationFactory<Program> factory)
    {
        _client = factory.CreateClient();
        _client.DefaultRequestHeaders.Add("X-API-Key", "dev-key");
    }

    [Fact]
    public async Task Idempotency_Returns_Same_Relay()
    {
        var body = new
        {
            eventType = "order.created",
            destination = new { type = "webhook", url = "https://example.com" },
            payload = new { x = 1 }
        };

        var req1 = await _client.PostAsJsonAsync("/v1/relays", body);
        var req2 = await _client.PostAsJsonAsync("/v1/relays", body);

        Assert.Equal(HttpStatusCode.Created, req1.StatusCode);
        Assert.Equal(HttpStatusCode.Created, req2.StatusCode);
    }
}
