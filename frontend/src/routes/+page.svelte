<script lang="ts">
import { onMount } from "svelte";

let quizzes: unknown[] = []; // Change type to unknown[]
let websocket: WebSocket;

async function getQuizzes() {
	try {
		const resp = await fetch("http://localhost:3000/api/quizzes", {
			method: "GET",
			headers: {
				"Content-Type": "application/json",
			},
		});

		console.log(resp);
		if (resp.ok) {
			quizzes = await resp.json();
		}
	} catch (err) {
		console.error(err);
	}
}

function connect() {
	websocket = new WebSocket("ws://localhost:3000/ws/1");

	websocket.onopen = () => {
		console.log("Connected to WebSocket");
	};

	websocket.onmessage = (event) => {
		console.log("Received message:", event.data);
		// Handle incoming WebSocket messages here
	};

	websocket.onclose = () => {
		console.log("WebSocket connection closed");
		// Optionally, you can attempt to reconnect here
	};

	websocket.onerror = (error) => {
		console.error("WebSocket error:", error);
	};
}

function sendMessage(message: string) {
	if (websocket.readyState === WebSocket.OPEN) {
		websocket.send(message);
	} else {
		console.error("WebSocket is not open");
	}
}

function onMessage(event: MessageEvent) {
	console.log("Received message:", event.data);
	// Handle incoming WebSocket messages here
	// For example, you can update the UI or store the received data
}

onMount(() => {
	connect();

	return () => {
		if (websocket) {
			websocket.close();
		}
	};
});
</script>

<div>
  <h1>Hello world!</h1>

  <button on:click={getQuizzes}>Get Quizzes</button>
  
  <pre>
    {#each quizzes as quiz}
      <pre>{JSON.stringify(quiz, null, 2)}</pre>
    {/each}
  </pre>

  <!-- Add a button to send a test socket message -->
  <button on:click={() => sendMessage('Test message')}>Send Test Message</button>

</div>
