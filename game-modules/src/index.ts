import { Repository } from "./repository/db";

async function main() {
	const repo = new Repository();
	try {
		const data = await repo.getCards();
		console.log(data);
		repo.close();
	} catch (error) {
		console.error(error);
		repo.close();
	}
}

main();
