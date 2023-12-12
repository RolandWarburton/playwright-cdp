import React, { useEffect, useState } from 'react';
import { Outlet } from 'react-router-dom';
import Navigation from '../../components/navigation';
import { v4 as uuidv4 } from 'uuid';

function setSessionID(sessionID: string) {
  localStorage.setItem('sessionID', sessionID);
}
function getSessionID() {
  const curSessionID = localStorage.getItem('sessionID');
  if (!curSessionID) {
    const sessionID = uuidv4();
    setSessionID(sessionID);
    return sessionID;
  } else {
    return curSessionID;
  }
}

function Root() {
  const [sessionID, setSessionID] = useState('');
  const [ws, setWS] = useState<WebSocket | undefined>(undefined);
  useEffect(() => {
    const sessionID = getSessionID();
    setSessionID(sessionID);
    const ws = new WebSocket(`ws://localhost:3000/ws/${sessionID}`);
    setWS(ws);

    ws.onopen = () => {
      console.log('WebSocket connection established');
      setSessionID(sessionID);
      ws.send('Hello WebSocket!');
    };

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log('Received message:', message);
    };

    ws.onclose = () => {
      console.log('WebSocket connection closed');
    };

    return () => {
      ws.close();
    };
  }, []);

  const sendPing = () => {
    if (!ws) {
      console.log('ws not established yet');
      return;
    }
    ws.send('ping');
  };

  return (
    <>
      <div>
        <Navigation />
        <h1>Root</h1>
        <button onClick={sendPing}>send ping</button>
        {sessionID}
        {/* render children here */}
        <Outlet />
      </div>
    </>
  );
}

export default Root;
