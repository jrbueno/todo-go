import {Heading, Stack} from "@chakra-ui/react";
import TodoItem from "./TodoItem.tsx";
import {useQuery} from "@tanstack/react-query";

export function TodoList() {
    const query = useQuery({ queryKey: ["todos"], 
        queryFn: async () => {
            try {
                const resp = await fetch("http://localhost:3000/api/todos");
                const data = await resp.json();
                if (!resp.ok) {
                    throw new Error("Network response was not ok\n" + data.error);
                }
                return data || [];
            } catch (err) {
                throw new Error(err as string);
            }
        }
    });
    return (
        <>
            <Heading bgGradient='linear(to-l, #7928CA, #FF0080)'
                     bgClip='text'
                     fontSize='6xl'
                     fontWeight='extrabold'>Todo List</Heading>
            <Stack gap={"3"}>
                {query.data?.map((todo, index) => {
                    return (
                        <TodoItem key={todo.id} todo={todo}/>
                    )
                })}
            </Stack>
        </>
    )
}