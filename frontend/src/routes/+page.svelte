<script lang="ts">
import QuizCard from "$lib/components/QuizCard.svelte";
import { onMount } from "svelte";

let quizzes: { id: string; name: string }[] = [];
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
	websocket = new WebSocket("ws://localhost:3000/ws");

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
  <h1 class="text-4xl font-bold">Hello world!</h1>


  <!-- Add a button to send a test socket message -->
  <div>
  	<button on:click={() => sendMessage('Test message')}>Send Test Message</button>
  </div>

  <div>
    <button on:click={getQuizzes}>Get Quizzes</button>
  </div>
  
  <ul class="flex flex-col gap-y-2">
    {#each quizzes as quiz}
      <li>
        <QuizCard name={quiz.name} id={quiz.id} />
      </li>
    {/each}
  </ul>
</div>
