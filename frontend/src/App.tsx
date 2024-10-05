import { useEffect, useState } from "react";
import "./App.css";
import AuthForms from "./components/custom/Form/AuthForm";
import Navbar from "./components/custom/Navbar";
import PokemonCollection from "./components/custom/PokemonCollection";

function App() {
	const [isAuthenticated, setIsAuthenticated] = useState(false);
	useEffect(() => {
		const token = localStorage.getItem("token");
		if (token) {
			setIsAuthenticated(true);
		}
	}, []);
	return (
		<>
			{isAuthenticated ? (
				<>
					<Navbar />
					<PokemonCollection />
				</>
			) : (
				<AuthForms />
			)}
		</>
	);
}

export default App;
