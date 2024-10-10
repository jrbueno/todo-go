import './App.css'
import {Container, Stack} from "@chakra-ui/react";
import {NavBar} from "./components/NavBar.tsx";
import AddNewTodo from "./components/AddNewTodo.tsx";
import {TodoList} from "./components/TodoList.tsx";

function App() {

  return (
    <Stack h="100vh">
        <NavBar />
        <Container>
            <AddNewTodo />
            <TodoList />
        </Container>
    </Stack>
  )
}

export default App
