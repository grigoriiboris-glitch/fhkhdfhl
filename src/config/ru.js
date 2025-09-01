// Шрифты
export const fontFamilyList = [
  { name: 'Song Ti', value: '宋体, SimSun, Songti SC' },
  { name: 'Microsoft Yahei', value: '微软雅黑, Microsoft YaHei' },
  { name: 'Italics', value: '楷体, 楷体_GB2312, SimKai, STKaiti' },
  { name: 'Boldface', value: '黑体, SimHei, Heiti SC' },
  { name: 'Official script', value: '隶书, SimLi' },
  { name: 'Andale Mono', value: 'andale mono' },
  { name: 'Arial', value: 'arial, helvetica, sans-serif' },
  { name: 'arialBlack', value: 'arial black, avant garde' },
  { name: 'Comic Sans Ms', value: 'comic sans ms' },
  { name: 'Impact', value: 'impact, chicago' },
  { name: 'Times New Roman', value: 'times new roman' },
  { name: 'Sans-Serif', value: 'sans-serif' },
  { name: 'serif', value: 'serif' }
]

// Стиль границы
export const borderDasharrayList = [
  { name: 'Сплошная', value: 'none' },
  { name: 'Штрих 1', value: '5,5' },
  { name: 'Штрих 2', value: '10,10' },
  { name: 'Штрих 3', value: '20,10,5,5,5,10' },
  { name: 'Штрих 4', value: '5, 5, 1, 5' },
  { name: 'Штрих 5', value: '15, 10, 5, 10, 15' },
  { name: 'Штрих 6', value: '1, 5' }
]

// Линии
export const lineStyleList = [
  { name: 'Прямая', value: 'straight' },
  { name: 'Кривая', value: 'curve' },
  { name: 'Напрямую', value: 'direct' }
]

// Корневая линия в «Кривая»
export const rootLineKeepSameInCurveList = [
  { name: 'Скобка', value: false },
  { name: 'Фигурная скобка', value: true }
]

// Повтор фона
export const backgroundRepeatList = [
  { name: 'Без повтора', value: 'no-repeat' },
  { name: 'Повтор', value: 'repeat' },
  { name: 'По X', value: 'repeat-x' },
  { name: 'По Y', value: 'repeat-y' }
]

// Позиция фона
export const backgroundPositionList = [
  { name: 'По умолчанию', value: '0% 0%' },
  { name: 'Слева сверху', value: 'left top' },
  { name: 'Слева по центру', value: 'left center' },
  { name: 'Слева снизу', value: 'left bottom' },
  { name: 'Справа сверху', value: 'right top' },
  { name: 'Справа по центру', value: 'right center' },
  { name: 'Справа снизу', value: 'right bottom' },
  { name: 'Сверху по центру', value: 'center top' },
  { name: 'По центру', value: 'center center' },
  { name: 'Снизу по центру', value: 'center bottom' }
]

// Размер фона
export const backgroundSizeList = [
  { name: 'Авто', value: 'auto' },
  { name: 'Cover', value: 'cover' },
  { name: 'Contain', value: 'contain' }
]

// Список форм
export const shapeList = [
  { name: 'Прямоугольник', value: 'rectangle' },
  { name: 'Ромб', value: 'diamond' },
  { name: 'Параллелограмм', value: 'parallelogram' },
  { name: 'Скругл. прямоугольник', value: 'roundedRectangle' },
  { name: 'Восьмиугольный прямоугольник', value: 'octagonalRectangle' },
  { name: 'Наружный треуг. прямоугольник', value: 'outerTriangularRectangle' },
  { name: 'Внутр. треуг. прямоугольник', value: 'innerTriangularRectangle' },
  { name: 'Эллипс', value: 'ellipse' },
  { name: 'Круг', value: 'circle' }
]

// Боковая панель
export const sidebarTriggerList = [
  { name: 'Стиль узла', value: 'nodeStyle', icon: 'iconzhuti' },
  { name: 'Базовый стиль', value: 'baseStyle', icon: 'iconyangshi' },
  { name: 'Тема', value: 'theme', icon: 'iconjingzi' },
  { name: 'Структура', value: 'structure', icon: 'iconjiegou' },
  { name: 'Контур', value: 'outline', icon: 'iconfuhao-dagangshu' },
  { name: 'Горячие клавиши', value: 'shortcutKey', icon: 'iconjianpan' }
]

// Типы скачивания
export const downTypeList = [
  { name: 'Специальный файл', type: 'smm', icon: 'iconwenjian', desc: 'Доступно для импорта' },
  { name: 'JSON', type: 'json', icon: 'iconjson', desc: 'Формат обмена данными, можно импортировать' },
  { name: 'Изображение', type: 'png', icon: 'iconPNG', desc: 'Подходит для просмотра и обмена' },
  { name: 'SVG', type: 'svg', icon: 'iconSVG', desc: 'Векторная графика' },
  { name: 'PDF', type: 'pdf', icon: 'iconpdf', desc: 'Подходит для печати' },
  { name: 'Markdown', type: 'md', icon: 'iconmarkdown', desc: 'Легко открыть другим ПО' },
  { name: 'XMind', type: 'xmind', icon: 'iconxmind', desc: 'Файл XMind' }
]

// Языки
export const langList = [
  { value: 'zh', name: '简体中文' },
  { value: 'en', name: 'English' },
  { value: 'ru', name: 'Русский' }
]

