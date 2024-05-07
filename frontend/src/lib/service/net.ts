export class NetService {
  private webSocket!: WebSocket;
  private textDecoder: TextDecoder = new TextDecoder();
  private textEncoder: TextEncoder = new TextEncoder();

  connect() {
    this.webSocket = new WebSocket("ws://localhost:3000/ws");

    this.webSocket.onopen = () => {
      console.log("Connected to server");
    };

    this.webSocket.onmessage = async (event: MessageEvent) => {
      const arrayBuffer = await event.data.arrayBuffer();
      const bytes = new Uint8Array(arrayBuffer);
      const packetId = bytes[0];

      const packet = JSON.parse(this.textDecoder.decode(bytes.subarray(1)));

      console.log(packet);
    };

    this.webSocket.onclose = () => {
      console.log("Disconnected from server");
    };
  }

  /**
   * Sends a packet to the server over the WebSocket connection.
   * @param packet The packet data to be sent.
   */
  sendPacket(packet: unknown) {
    // Define the packet ID (in this case, a fixed value of 1337)
    const packetId = 1337;

    // Convert the packet data to a JSON string
    const packetData = JSON.stringify(packet);

    // Create a Uint8Array to hold the packet ID (1 byte)
    const packetIdArray = new Uint8Array([packetId]);

    // Encode the packet data string into a Uint8Array
    const packetDataArray = this.textEncoder.encode(packetData);

    // Create a new Uint8Array to hold the merged packet data
    // Its size is the sum of the lengths of packetIdArray and packetDataArray
    const mergedArray = new Uint8Array(
      packetIdArray.length + packetDataArray.length,
    );

    // Copy the packet ID array into the merged array starting at index 0
    mergedArray.set(packetIdArray);

    // Copy the packet data array into the merged array starting after the packet ID
    mergedArray.set(packetDataArray, packetIdArray.length);

    // Send the merged array (packet ID + packet data) over the WebSocket connection
    this.webSocket.send(mergedArray);
  }
}
