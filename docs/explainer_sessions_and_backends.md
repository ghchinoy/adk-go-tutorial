# Explainer: Sessions and Backends

When building agents with ADK for Go, it's important to understand the distinction between the **Model Backend** (where the LLM runs) and the **Session Service** (where conversation history and state are stored).

## Model Backend

The Model Backend determines which service processes your prompts and generates responses. ADK supports both Google AI Studio and Vertex AI.

You configure this when creating the model:

```go
// Option 1: Google AI Studio (Default if Backend is omitted)
model, _ := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    APIKey: os.Getenv("GOOGLE_API_KEY"),
})

// Option 2: Vertex AI
model, _ := gemini.NewModel(ctx, "gemini-2.5-flash", &genai.ClientConfig{
    Backend:  genai.BackendVertexAI,
    Project:  os.Getenv("GOOGLE_CLOUD_PROJECT"),
    Location: os.Getenv("GOOGLE_CLOUD_LOCATION"),
})
```

Changing the backend *only* changes where the generation request is sent. It does not automatically change how conversation history is managed.

## Session Service

The Session Service is responsible for storing the state of the conversation. This includes:
*   **Message History**: The back-and-forth between the user and the agent.
*   **Agent State**: Variables and data that the agent needs to remember across turns.

ADK provides two built-in implementations and allows for custom ones.

### 1. In-Memory Session (Default)

If you do not specify a `SessionService` in your `adk.Config`, the launcher defaults to an **in-memory** service.

*   **Pros**: Zero configuration, extremely fast, ideal for local development, testing, and simple CLI tools.
*   **Cons**: **Ephemeral**. If your application restarts or crashes, all conversation history and state are lost instantly. Not suitable for production web services where users expect continuity.

### 2. Vertex AI Session Service (Managed)

The `session.VertexAIService` is backed by the **Vertex AI Agent Engine**.

*   **Pros**:
    *   **Fully Managed**: Google handles the storage, scaling, and security of the session data.
    *   **Integrated**: Designed to work seamlessly with other Vertex AI services.
*   **Cons**:
    *   Requires Google Cloud setup (project, billing, permissions).
    *   Less direct control over the underlying storage mechanism.

### 3. Custom Session Service (e.g., Firestore, Redis)

Because `session.Service` is a standard Go interface, you can implement your own. This is a powerful option for production applications with specific requirements.

*   **Pros**:
    *   **Total Control**: You decide exactly how and where data is stored (e.g., Firestore for document storage, Redis for high-speed caching, PostgreSQL for relational data).
    *   **Portable**: Your session data can live alongside your other application data.
    *   **Cost Optimization**: You can choose a storage solution that fits your specific budget and performance needs.
*   **Cons**:
    *   **Development Effort**: You must write and maintain the code to implement the `session.Service` interface (Create, Get, List, Delete, AppendEvent).

#### Implementing a Custom Service

To build a custom service, you just need to implement this interface:

```go
type Service interface {
	Create(context.Context, *CreateRequest) (*CreateResponse, error)
	Get(context.Context, *GetRequest) (*GetResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Delete(context.Context, *DeleteRequest) error
	AppendEvent(context.Context, Session, *Event) error
}
```

For example, a Firestore implementation would map `Create` to `firestoreClient.Collection("sessions").Add(...)` and `AppendEvent` to adding a new document to a "events" subcollection.

## Summary Comparison

| Feature | In-Memory | Vertex AI Agent Engine | Custom (e.g., Firestore) |
| :--- | :--- | :--- | :--- |
| **Persistence** | No (Lost on restart) | Yes (Managed by Google) | Yes (Managed by you) |
| **Setup Effort** | None | Low (GCP Config) | High (Code + Infra) |
| **Control** | N/A | Low (Black box) | High (Full schema control) |
| **Scalability** | Limited by RAM | High (Cloud scale) | Depends on chosen backend |
| **Best For** | Testing, CLIs | Rapid production deployment | Complex apps with specific data needs |