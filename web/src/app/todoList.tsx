'use client';
import useWebSocket from 'react-use-websocket';
import styles from './page.module.css';
import { useState } from 'react';

const socketUrl = 'ws://127.0.0.1:1234';

export default function TodoList() {
	const [list, setList] = useState<string[]>([]);
	const [input, setInput] = useState<string>('');
	const [listName, setListName] = useState<string>('test');
	const [userId, setUserId] = useState<string>('default');
	const [channelId, setChannelId] = useState<number>(0);
	const {
		sendMessage,
		sendJsonMessage,
		lastMessage,
		lastJsonMessage,
		readyState,
		getWebSocket,
	} = useWebSocket(socketUrl, {
		onOpen: () => {
			console.log('opened');
			sendJsonMessage({ userId: userId });
		},
		onClose: (e) => console.log('closed: ' + e.reason),
		onMessage: (e) => {
			console.log(e.data);
			const data: {
				channelId: number;
				msg: string;
			} = JSON.parse(e.data);
			setChannelId(data.channelId);
			if (data.msg === '') {
				// get list http request
				fetch('http://127.0.0.1:8080/call/get-todo-list', {
					method: 'POST',
					headers: {
						'Content-Type': 'application/json',
					},
					body: JSON.stringify({
						name: listName,
						channelId: data.channelId,
						userId: userId,
					}),
				})
					.then((response) => response.json())
					.then((data) => {
						console.log('Success:', data);
					})
					.catch((error) => {
						console.error('Error:', error);
					});
				return;
			} else {
				alert(data.msg);
			}
		},
		//Will attempt to reconnect on all close events, such as server shutting down
		shouldReconnect: (closeEvent) => true,
	});

	return (
		<>
			<input
				value={listName}
				placeholder='list name'
				onChange={(e) => setListName(e.target.value)}
			></input>
			<input
				value={userId}
				placeholder='user id'
				onChange={(e) => setUserId(e.target.value)}
			></input>
			<div className={styles.description}>
				<input
					value={input}
					placeholder='item content'
					onChange={(e) => setInput(e.target.value)}
				></input>
				<button
					onClick={async () => {
						// create item http request
						if (channelId > 0) {
							fetch('http://127.0.0.1:8080/call/add-todo-item', {
								method: 'POST',
								headers: {
									'Content-Type': 'application/json',
								},
								body: JSON.stringify({
									name: listName,
									channelId: channelId,
									userId: userId,
									description: input,
									completed: false,
								}),
							})
								.then((response) => response.json())
								.then((data) => {
									console.log('Success:', data);
								})
								.catch((error) => {
									console.error('Error:', error);
								});
						}
					}}
				>
					clink me to add item
				</button>
				<div>
					{list.map((v, i) => {
						return <p key={i}>{v}</p>;
					})}
				</div>
			</div>
		</>
	);
}
