import {Badge, Box, Flex, Spinner, Text} from "@chakra-ui/react";
import {FaCheckCircle} from "react-icons/fa";
import {MdDelete} from "react-icons/md";
import {useMutation, useQueryClient} from "@tanstack/react-query";

export type Todo = {
    id: number;
    title: string;
    completed: boolean;
}

const TodoItem = ({todo}: { todo: Todo }) => {
    const queryClient = useQueryClient();
    // Update Mutations
    const {mutate: updateTodo, isPending: isUpdating} = useMutation({
        mutationKey: ["updateTodo"],
        mutationFn: async () => {
            if (todo.completed) return alert("Todo is already completed");
            const resp = await fetch("http://localhost:3000/api/todos/" + todo.id, {
                method: "PATCH",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({completed: true}),
            });
            const data = await resp.json();
            if (!resp.ok) {
                throw new Error("Network response was not ok\n" + data.error);
            }
            return data;
        },
        onSuccess: () => {
            // Invalidate and fetch
            queryClient.invalidateQueries({ queryKey: ['todos'] })
        },
    })
    //Delete Mutations
    const {mutate: deleteTodo, isPending: isDeleting} = useMutation({
        mutationKey: ["deleteTodo"],
        mutationFn: async () => {
            const resp = await fetch("http://localhost:3000/api/todos/" + todo.id, {
                method: "DELETE",
            });
            if (!resp.status === 204) {
                const data = await resp.json();
                throw new Error("Network response was not ok\n" + data.error);
            }
            return;
        },
        onSuccess: () => {
            // Invalidate and fetch
            queryClient.invalidateQueries({ queryKey: ['todos'] })
        },
        onError: (error) => {
            alert("An error occurred: " + error);
        }
    })
    return (
        <Flex gap={2} alignItems={"center"}>
            <Flex
                flex={1}
                alignItems={"center"}
                border={"1px"}
                borderColor={"gray.600"}
                p={2}
                borderRadius={"lg"}
                justifyContent={"space-between"}
            >
                <Text
                    color={todo.completed ? "green.200" : "black.100"}
                    textDecoration={todo.completed ? "line-through" : "none"}
                >
                    {todo.title}
                </Text>
                {todo.completed && (
                    <Badge ml='1' colorScheme='green'>
                        Done
                    </Badge>
                )}
                {!todo.completed && (
                    <Badge ml='1' colorScheme='yellow'>
                        In Progress
                    </Badge>
                )}
            </Flex>
            <Flex gap={2} alignItems={"center"}>
                <Box color={"green.500"} cursor={"pointer"} onClick={updateTodo}>
                    {!isUpdating && <FaCheckCircle size={20} />}
                    {isUpdating && <Spinner size={"sm"} />}
                </Box>
                <Box color={"red.500"} cursor={"pointer"} onClick={deleteTodo}>
                    { !isDeleting && <MdDelete size={25} />}
                    { isDeleting && <Spinner size={"sm"} />}
                </Box>
            </Flex>
        </Flex>
    );
};
export default TodoItem;