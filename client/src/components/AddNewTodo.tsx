import { Button, Flex, Input, Spinner } from "@chakra-ui/react";
import { useState } from "react";
import { IoMdAdd } from "react-icons/io";
import {useMutation, useQueryClient} from "@tanstack/react-query";
import {Todo} from "./TodoItem.tsx";
import * as React from "react";

const AddNewTodo = () => {
	const [newTodo, setNewTodo] = useState("");
	const queryClient = useQueryClient();
	const { mutate: createTodo, isPending } = useMutation<Todo>({
		mutationKey: ["createTodo"],
		mutationFn: async (e: React.FormEvent) => {
			e.preventDefault();
			const resp = await fetch("http://localhost:3000/api/todos", {
				method: "POST",
				headers: {
					"Content-Type": "application/json",
				},
				body: JSON.stringify({ title: newTodo}),
			});
			const data = await resp.json();
			if (!resp.ok) {
				throw new Error("Network response was not ok\n" + data.error);
			}
			return data;
		},
		onSuccess: () => {
			queryClient.invalidateQueries({ queryKey: ["todos"] });
			setNewTodo("");
		},
		onError: (error) => {
			alert("An error occurred: " + error);
		},
	});
	return (
		<form onSubmit={createTodo}>
			<Flex gap={2}>
				<Input
					type='text'
					value={newTodo}
					onChange={(e) => setNewTodo(e.target.value)}
					ref={(input) => input && input.focus()}
				/>
				<Button
					mx={2}
					type='submit'
					_active={{
						transform: "scale(.97)",
					}}
				>
					{isPending ? <Spinner size={"xs"} /> : <IoMdAdd size={30} />}
				</Button>
			</Flex>
		</form>
	);
};
export default AddNewTodo;