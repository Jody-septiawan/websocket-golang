// import useContext hook
import React, { useState, useEffect, useContext } from "react";

import Navbar from "../components/Navbar";

import { Container, Row, Col } from "react-bootstrap";
import Contact from "../components/complain/Contact";
// import chat component
import Chat from "../components/complain/Chat";

// import user context
import { UserContext } from "../context/userContext";

// import socket.io-client
import { io } from "socket.io-client";

import * as WebSocket from "websocket";

// initial variable outside socket
let socket;
export default function Complain() {
  const [contact, setContact] = useState(null);
  const [contacts, setContacts] = useState([]);
  // create messages state
  const [messages, setMessages] = useState([]);

  const title = "Complain admin";
  document.title = "DumbMerch | " + title;

  // consume user context
  const [state] = useContext(UserContext);

  useEffect(() => {
    socket = new WebSocket.w3cwebsocket("ws://localhost:8080/ws");

    socket.onopen = function () {
      console.log("Client connect to Server Socket");

      loadContact();

      socket.onmessage = (res) => {
        res = JSON.parse(res.data);
        if (res.event == "admin contact") {
          const data = res.data.map((item) => ({
            ...item,
            message:
              item.chats.length > 0
                ? item.chats[item.chats.length - 1].message
                : "Click here to start message",
          }));
          setContacts(data);
        } else if (res.event == "messages") {
          setMessages(res.data);
        } else if (res.event == "new message") {
          loadMessages();
          //   loadContact();
        }
      };

      if (contact != null) {
        loadMessages();
      }
    };
  }, [contact]);

  const loadContact = () => {
    socket.send(
      JSON.stringify({
        event: "admin contact",
        senderId: state.user.id,
      })
    );
  };

  // used for active style when click contact
  const onClickContact = (data) => {
    setContact(data);
  };

  const loadMessages = (value) => {
    console.log("contact?.id : ", contact);
    console.log("state : ", state);

    socket.send(
      JSON.stringify({
        event: "messages",
        senderId: state.user.id,
        recipientId: contact?.id,
      })
    );
  };
  const onSendMessage = (e) => {
    // listen only enter key event press
    if (e.key === "Enter") {
      //emit event send message
      socket.send(
        JSON.stringify({
          event: "send message",
          senderId: state.user.id,
          recipientId: contact.id,
          message: e.target.value,
        })
      );
      e.target.value = "";
    }
  };

  return (
    <>
      <Navbar title={title} />
      <Container fluid style={{ height: "89.5vh" }}>
        <Row>
          <Col
            md={3}
            style={{ height: "89.5vh" }}
            className="px-3 border-end border-dark overflow-auto"
          >
            <Contact
              dataContact={contacts}
              clickContact={onClickContact}
              contact={contact}
            />
          </Col>
          <Col md={9} style={{ maxHeight: "89.5vh" }} className="px-0">
            <Chat
              contact={contact}
              messages={messages}
              user={state.user}
              sendMessage={onSendMessage}
            />
          </Col>
        </Row>
      </Container>
    </>
  );
}
