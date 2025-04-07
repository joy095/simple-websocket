import React, { useEffect, useRef, useState } from "react";

const GoChat: React.FC = () => {
  const [username, setUsername] = useState("");
  const [message, setMessage] = useState("");
  const [messages, setMessages] = useState<{ from: string; message: string }[]>(
    []
  );
  const ws = useRef<WebSocket | null>(null);

  useEffect(() => {
    ws.current = new WebSocket("ws://localhost:8080/ws");

    ws.current.onmessage = (event) => {
      const msg = JSON.parse(event.data);
      setMessages((prev) => [...prev, msg]);
    };

    return () => {
      ws.current?.close();
    };
  }, []);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    if (ws.current?.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify({ from: username, message }));
      setMessage(""); // Clear input after sending
    }
  };

  return (
    <div style={{ display: "flex", flexDirection: "column" }}>
      <div id="inputs">
        <form onSubmit={handleSubmit}>
          <input
            type="text"
            name="username"
            placeholder="username"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            required
          />
          <input
            type="text"
            name="message"
            placeholder="what's on your mind ?"
            value={message}
            onChange={(e) => setMessage(e.target.value)}
            required
          />
          <input type="submit" value="Send" />
        </form>
      </div>
      <div id="messages" style={{ display: "flex", flexDirection: "column" }}>
        {messages.map((msg, idx) => (
          <p key={idx}>
            <b>{msg.from}</b> says {msg.message}
          </p>
        ))}
      </div>
    </div>
  );
};

export default GoChat;
