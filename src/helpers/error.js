const map = {
	phone: 'Телефон',
	email: 'Емайл',
	password: 'Пароль',
	hash: 'Капча'
}

export default {
	getError(e) {
		let text = '';
		if (e.response) {
			for (let [key, value] of Object.entries(e.response.data.errors)) {
				value = value.join(',')
				if (key in map) {
					text += map[key] + ' ' + value.replace(':', '');
				}
				
			}
		}

		return text || 'Что-то пошло не так :(';
	}
}