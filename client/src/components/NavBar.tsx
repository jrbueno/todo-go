import {IconButton, Stack, useColorMode, useColorModeValue} from "@chakra-ui/react";
import reactLogo from "../assets/react.svg";
import {IoMoon} from "react-icons/io5";
import {LuSun} from "react-icons/lu";

// import viteLogo from "*.svg";

export function NavBar() {
    const { colorMode, toggleColorMode } = useColorMode()
    //Create a navbar component that is centered and has icons on the left and a dark mode toggle on the right
    return (
        <Stack w="800px" direction="row" justify="space-between" align="center" p={4}
               bg={useColorModeValue("gray.400", "gray.700")} borderRadius={"5"}>
            <Stack direction="row" align="center">
                <img src={reactLogo} alt="React Logo" className="App-logo"/>
                {/*<img src={viteLogo} alt="Vite Logo" className="App-logo"/>*/}
            </Stack>
            <IconButton onClick={toggleColorMode} aria-label="Toggle Dark Mode"
                        icon={colorMode === "light" ? <IoMoon/> : <LuSun size={20}/>}/>
        </Stack>
    )
}