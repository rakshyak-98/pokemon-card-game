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
				<div className="flex min-h-screen bg-background text-text-primary">
					<Navbar />
					<PokemonCollection />
				</div>
			) : (
				<AuthForms />
			)}
		</>
	);
}

export default App;
