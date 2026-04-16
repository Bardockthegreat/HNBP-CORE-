# Example: AI Node

## Scenario

An AI agent (`gpt-agent-v2`) is performing a legal document summarization task. Every action — initialization, tool calls, and output production — is logged to an HNBP chain. This creates a verifiable audit trail of the AI's reasoning and actions.

## How HNBP-CORE Is Used

1. The AI agent initializes its log on task start (genesis record).
2. Every tool call, API call, or decision is appended as a record, with relevant metadata (tool name, args hash, result hash).
3. The final output hash is appended so any downstream consumer can verify the AI produced exactly what the log says.
4. The log is exposed to the user and any auditor.

## Key Benefit

This makes AI behavior **auditable and disputable**. If a user claims the AI did something wrong, the log proves exactly what happened. If the AI is accused of hallucinating, the tool call log shows what data it actually accessed.

## Use in Production

```python
from hnbp import append, validate_log, save_log

log = []
append(log, {"actor": "my-agent", "event": "init"})
# ... on each action:
append(log, {"actor": "my-agent", "event": "tool-call", "tool": "search", "query": "..."})
save_log(log, "agent.morus")
```
