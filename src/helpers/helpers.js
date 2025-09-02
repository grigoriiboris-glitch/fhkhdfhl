const DateOptions = {
  year: 'numeric',
  month: 'long',
  day: 'numeric',
  hour: '2-digit',
  minute: '2-digit'
};

export default {
  moveElement(arr, fromIndex, toIndex) {
    if (toIndex >= arr.length || toIndex < 0) {
      console.log("Некорректный индекс назначения");
      return arr;
    }
  
    // Удаляем элемент с позиции fromIndex и сохраняем его
    const element = arr.splice(fromIndex, 1)[0];
  
    // Вставляем элемент на новую позицию
    arr.splice(toIndex, 0, element);
  
    return arr;
  },
  toCamelCase(str) {
    let name = str.split('_');
    let nameUp = '';
    name.forEach(el => {
      nameUp += el.charAt(0).toUpperCase() + el.slice(1);
    })
    return nameUp;
  },
  proxyToArray(proxyObject) {
    let arr = [];
    proxyObject.forEach(el => {
      if (Array.isArray(el)) {
        arr.push(this.proxyToArray(el));

      } else if (typeof el == 'object') {
        arr.push(this.proxyToObj(el));

      } else {
        arr.push(el);
      }
    })

    return arr;
  },
  proxyToObj(obj1) {
    let newObj = {};
    for (let attrname in obj1) {
      if (Array.isArray(obj1[attrname])) {
        newObj[attrname] = this.proxyToArray(obj1[attrname])

        continue;
      }
      if (typeof obj1[attrname] == 'object') {
        newObj[attrname] = this.proxyToObj(obj1[attrname])
      } else {
        newObj[attrname] = obj1[attrname];
      }
    }
    return newObj;
  },
  mergeObj(arr) {
    let newObj = {};

    arr.forEach((obj1) => {
      for (let attrname in obj1) {
        if (Array.isArray(obj1[attrname])) {
          newObj[attrname] = this.proxyToArray(obj1[attrname]);
          continue;
        }

        if (typeof obj1[attrname] == 'object') {
          newObj[attrname] = this.proxyToObj(obj1[attrname]);
          continue;
        } else {
          newObj[attrname] = obj1[attrname];
        }
      }
    });
    return newObj;
  },
  openTarget(url) {
    let x = screen.width / 2 - 700 / 2;
    let y = screen.height / 2 - 450 / 2;
    window.open(url, '_blank', `location=yes,height=570,width=855,left=${x},top=${y}`);
  },
  addUrlParams(url, params) {
    // Check if the URL already contains a query string
    const separator = url.includes('?') ? '&' : '?';

    // Create an array of key-value pairs for the parameters
    const paramPairs = Object.entries(params);

    // Concatenate the parameters to the URL
    const paramString = paramPairs.map(pair => pair.join('=')).join('&');

    // Return the new URL with added parameters
    return url + separator + paramString;
  },
  formatDateForHuman(date, locale = 'ru') {
    var today = new Date(date);
    // en-US
    date = today.toLocaleDateString(locale, DateOptions);

    return date;
  },
  clean(data) {
    for (let k in data) {
      data[k] = null;
    }
    log(data)
    return data;
  },
  cliceWord(val, length = 30) {
    if (val === null) {
      return '';
    }
    if (val.length >= length) {
      val = val.slice(0, length);
      val += '...';
    }
    return val;
  },
  getTime(date) {
    date = new Date(date);
    let hours = date.getHours().toString().padStart(2, '0');
    let minutes = date.getMinutes().toString().padStart(2, '0');

    return `${hours}:${minutes}`;
  },
  dateConvert(str, slash = '-') {
    if (str == '' || str === null) {
      return '';
    }
    var date = new Date(str),
      month = ('0' + (date.getMonth() + 1)).slice(-2),
      day = ('0' + date.getDate()).slice(-2);

    let year = date.getFullYear();

    return [day, month, year].join(slash);
  },
  boolToHuman(str) {
    return str === true ? 'да' : 'нет';
  },
  dateFormat(dateString, slash = '.') {
    // return new Date(dateString).toLocaleDateString('ru');

    /*if (dateString == null) {
return '';
}*/
    // safari fix
    const match = dateString.match(/\d{4}-\d{2}-\d{2}/) || [];
    dateString = match[0] || '';
    if (!dateString) return '';
    const date = new Date(dateString);
    const yyyy = date.getFullYear();
    let mm = date.getMonth() + 1;
    mm = mm.toString();

    let dd = date.getDate();
    dd = dd.toString();

    if (mm.length < 2) {
      mm = 0 + mm;
    }
    if (dd.length < 2) {
      dd = 0 + dd;
    }

    return `${dd}${slash}${mm}${slash}${yyyy}`;
  },
  downloadFile(data, filename, type) {
    var file = new Blob([data], { type: type });
    if (window.navigator.msSaveOrOpenBlob)
      // IE10+
      window.navigator.msSaveOrOpenBlob(file, filename);
    else {
      // Others
      var a = document.createElement('a'),
        url = URL.createObjectURL(file);
      a.href = url;
      a.download = filename;
      document.body.appendChild(a);
      a.click();
      setTimeout(function () {
        document.body.removeChild(a);
        window.URL.revokeObjectURL(url);
      }, 0);
    }
  }
};
