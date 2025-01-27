export function createStreamChat(request) {
  return new Promise((resolve, reject) => {
    const ws = new WebSocket('ws://localhost:8080/ws/chat');
    let fullContent = '';

    ws.onopen = () => {
      ws.send(JSON.stringify(request));
    };

    ws.onmessage = (event) => {
      const response = JSON.parse(event.data);
      
      if (response.error) {
        ws.close();
        reject(new Error(response.error));
        return;
      }

      if (response.done) {
        ws.close();
        resolve(fullContent);
        return;
      }

      fullContent += response.content;
      // 触发进度更新回调
      if (request.onProgress) {
        request.onProgress(fullContent);
      }
    };

    ws.onerror = (error) => {
      reject(error);
    };
  });
} 