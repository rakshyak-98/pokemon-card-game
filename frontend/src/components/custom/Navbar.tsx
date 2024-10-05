import { Button } from "@/components/ui/button";
import {
  ChevronLeft,
  ChevronRight,
  Home,
  List,
  Settings,
  User,
} from "lucide-react";
import React, { useState } from "react";

interface NavItem {
	name: string;
	icon: React.ReactNode;
	href: string;
}

const navItems: NavItem[] = [
	{ name: "Home", icon: <Home className="h-5 w-5" />, href: "/" },
	{
		name: "Collection",
		icon: <List className="h-5 w-5" />,
		href: "/collection",
	},
	{ name: "Profile", icon: <User className="h-5 w-5" />, href: "/profile" },
	{
		name: "Settings",
		icon: <Settings className="h-5 w-5" />,
		href: "/settings",
	},
];

export default function CollapsibleSidebar() {
	const [isExpanded, setIsExpanded] = useState(true);

	const toggleSidebar = () => {
		setIsExpanded(!isExpanded);
	};

	return (
		<div
			className={`fixed left-0 top-0 h-full bg-card text-text-primary transition-all duration-300 ${
				isExpanded ? "w-64" : "w-20"
			}`}
		>
			<div className="flex h-full flex-col justify-between p-4">
				<div>
					<div className="mb-8 flex items-center justify-between">
						{isExpanded && (
							<h2 className="text-xl font-bold text-accent-primary">Pokémon</h2>
						)}
						<Button
							variant="ghost"
							size="icon"
							onClick={toggleSidebar}
							aria-label={isExpanded ? "Collapse sidebar" : "Expand sidebar"}
							className="ml-auto"
						>
							{isExpanded ? (
								<ChevronLeft className="h-5 w-5" />
							) : (
								<ChevronRight className="h-5 w-5" />
							)}
						</Button>
					</div>
					<nav>
						<ul className="space-y-2">
							{navItems.map((item) => (
								<li key={item.name}>
									<a
										href={item.href}
										className={`flex items-center rounded-lg p-2 text-text-secondary hover:bg-card-hover hover:text-accent-primary ${
											isExpanded ? "justify-start space-x-4" : "justify-center"
										}`}
									>
										{item.icon}
										{isExpanded && <span>{item.name}</span>}
									</a>
								</li>
							))}
						</ul>
					</nav>
				</div>
				<div className="text-center text-sm text-text-secondary">
					{isExpanded && <p>&copy; 2024 PokéDex</p>}
				</div>
			</div>
		</div>
	);
}
