import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { X } from "lucide-react";
export default function PokemonDialog({
	selectedPokemon,
	closeModal,
}: {
	selectedPokemon: any;
	closeModal: any;
}) {
	return (
		<Dialog open={selectedPokemon !== null} onOpenChange={closeModal}>
			<DialogContent className="bg-card text-text-primary">
				<DialogHeader>
					<DialogTitle className="text-2xl font-bold text-accent-primary">
						{selectedPokemon?.name}
					</DialogTitle>
					<Button
						variant="ghost"
						size="icon"
						onClick={closeModal}
						className="absolute right-4 top-4"
						aria-label="Close"
					>
						<X className="h-4 w-4" />
					</Button>
				</DialogHeader>
				{selectedPokemon && (
					<div className="mt-4">
						<div className="bg-background rounded-lg p-4 mb-4">
							<img
								src={selectedPokemon.image}
								alt={selectedPokemon.name}
								className="w-full h-60 object-contain"
							/>
						</div>
						<p className="mb-2">{selectedPokemon.description}</p>
						<div className="grid grid-cols-2 gap-4">
							<div>
								<p className="text-text-secondary">
									Type:{" "}
									<span className="text-accent-primary">
										{selectedPokemon.type}
									</span>
								</p>
								<p className="text-text-secondary">
									HP:{" "}
									<span className="text-accent-primary">
										{selectedPokemon.hp}
									</span>
								</p>
								<p className="text-text-secondary">
									Height:{" "}
									<span className="text-accent-primary">
										{selectedPokemon.height} m
									</span>
								</p>
								<p className="text-text-secondary">
									Weight:{" "}
									<span className="text-accent-primary">
										{selectedPokemon.weight} kg
									</span>
								</p>
							</div>
							<div>
								<p className="text-text-secondary">Abilities:</p>
								<ul className="list-disc list-inside">
									{selectedPokemon.abilities.map((ability: any, index: any) => (
										<li key={index} className="text-accent-primary">
											{ability}
										</li>
									))}
								</ul>
							</div>
						</div>
					</div>
				)}
			</DialogContent>
		</Dialog>
	);
}
