export default class User {
	constructor(object) {

		for (const [key, value] of Object.entries(object)) {
			if (key !== 'bots') {
				this[key] = value
			}
		}

	}

	getMyRoles() {
		return this.roles.filter((el)=>el.owner_id != this.id || (el.owner_id === this.id && el.user_id === this.id));
	}

	getLogins(system) {
		let items = [];
		let data = this.socialite_login.filter((el) => el.provider === system);

		data.forEach(el => {
			items.push({
				name: el.data.email || 'account',
				value: el.id,
				img: el.data.avatar || null
			})
		})
		items.reverse();

		return items;
	}

	getProvidedRoles() {
		return this.roles.filter((el)=>el.owner_id === this.id && el.user_id !== this.id);
	}

	isSuperAdmin() {
		return this.role === 'admin';
	}

	isAdmin() {
		return this.role_id === 1;
	}

	isModerator() {
		return this.role_id === 2;
	}

	isManager() {
		return this.role_id === 3;
	}
}