import { Button } from "@/components/ui/button";
import {
  Card,
  CardContent,
  CardDescription,
  CardFooter,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { EyeIcon, EyeOffIcon } from "lucide-react";
import React, { useState } from "react";

export default function AuthForms() {
	const [showPassword, setShowPassword] = useState(false);
	const [activeTab, setActiveTab] = useState("login");

	const togglePasswordVisibility = () => setShowPassword(!showPassword);

	const handleSubmit = (event: React.FormEvent<HTMLFormElement>) => {
		event.preventDefault();
		// Here you would typically handle form submission, e.g., sending data to an API
		console.log(`${activeTab} form submitted`);
	};

	return (
		<Card className="w-full max-w-md mx-auto">
			<CardHeader>
				<CardTitle className="text-2xl font-bold text-center">
					Welcome to PokéDex
				</CardTitle>
				<CardDescription className="text-center">
					Sign in to your account or create a new one
				</CardDescription>
			</CardHeader>
			<CardContent>
				<Tabs value={activeTab} onValueChange={setActiveTab}>
					<TabsList className="grid w-full grid-cols-2">
						<TabsTrigger value="login">Login</TabsTrigger>
						<TabsTrigger value="signup">Sign Up</TabsTrigger>
					</TabsList>
					<TabsContent value="login">
						<form onSubmit={handleSubmit}>
							<div className="space-y-4">
								<div className="space-y-2">
									<Label htmlFor="login-email">Email</Label>
									<Input
										id="login-email"
										type="email"
										placeholder="johndoe@example.com"
										required
									/>
								</div>
								<div className="space-y-2">
									<Label htmlFor="login-password">Password</Label>
									<div className="relative">
										<Input
											id="login-password"
											type={showPassword ? "text" : "password"}
											placeholder="••••••••"
											required
										/>
										<Button
											type="button"
											variant="ghost"
											size="icon"
											className="absolute right-2 top-1/2 -translate-y-1/2"
											onClick={togglePasswordVisibility}
										>
											{showPassword ? (
												<EyeOffIcon className="h-4 w-4" />
											) : (
												<EyeIcon className="h-4 w-4" />
											)}
										</Button>
									</div>
								</div>
								<Button type="submit" className="w-full">
									Login
								</Button>
							</div>
						</form>
					</TabsContent>
					<TabsContent value="signup">
						<form onSubmit={handleSubmit}>
							<div className="space-y-4">
								<div className="space-y-2">
									<Label htmlFor="signup-name">Full Name</Label>
									<Input
										id="signup-name"
										type="text"
										placeholder="John Doe"
										required
									/>
								</div>
								<div className="space-y-2">
									<Label htmlFor="signup-email">Email</Label>
									<Input
										id="signup-email"
										type="email"
										placeholder="johndoe@example.com"
										required
									/>
								</div>
								<div className="space-y-2">
									<Label htmlFor="signup-password">Password</Label>
									<div className="relative">
										<Input
											id="signup-password"
											type={showPassword ? "text" : "password"}
											placeholder="••••••••"
											required
										/>
										<Button
											type="button"
											variant="ghost"
											size="icon"
											className="absolute right-2 top-1/2 -translate-y-1/2"
											onClick={togglePasswordVisibility}
										>
											{showPassword ? (
												<EyeOffIcon className="h-4 w-4" />
											) : (
												<EyeIcon className="h-4 w-4" />
											)}
										</Button>
									</div>
								</div>
								<Button type="submit" className="w-full">
									Sign Up
								</Button>
							</div>
						</form>
					</TabsContent>
				</Tabs>
			</CardContent>
			<CardFooter className="flex justify-center">
				<p className="text-sm text-text-secondary">
					By signing in, you agree to our{" "}
					<a href="#" className="text-accent-primary hover:underline">
						Terms of Service
					</a>{" "}
					and{" "}
					<a href="#" className="text-accent-primary hover:underline">
						Privacy Policy
					</a>
					.
				</p>
			</CardFooter>
		</Card>
	);
}
