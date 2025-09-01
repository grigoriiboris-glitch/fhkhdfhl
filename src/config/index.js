import {
  fontSizeList,
  lineHeightList,
  colorList,
  borderWidthList,
  borderRadiusList,
  lineWidthList,
  store,
  langList,
  fontFamilyList as fontFamilyListZh,
  borderDasharrayList as borderDasharrayListZh,
  lineStyleList as lineStyleListZh,
  rootLineKeepSameInCurveList as rootLineKeepSameInCurveListZh,
  backgroundRepeatList as backgroundRepeatListZh,
  backgroundPositionList as backgroundPositionListZh,
  shortcutKeyList as shortcutKeyListZh,
  shapeList as shapeListZh,
  sidebarTriggerList as sidebarTriggerListZh,
  backgroundSizeList as backgroundSizeListZh,
  downTypeList as downTypeListZh,
  shapeListMap as shapeListMapZh,
  lineStyleMap as lineStyleMapZh
} from './zh'

import {
  fontFamilyList as fontFamilyListEn,
  borderDasharrayList as borderDasharrayListEn,
  lineStyleList as lineStyleListEn,
  rootLineKeepSameInCurveList as rootLineKeepSameInCurveListEn,
  backgroundRepeatList as backgroundRepeatListEn,
  backgroundPositionList as backgroundPositionListEn,
  shortcutKeyList as shortcutKeyListEn,
  shapeList as shapeListEn,
  sidebarTriggerList as sidebarTriggerListEn,
  backgroundSizeList as backgroundSizeListEn,
  downTypeList as downTypeListEn
} from './en'

import {
  fontFamilyList as fontFamilyListRu,
  borderDasharrayList as borderDasharrayListRu,
  lineStyleList as lineStyleListRu,
  rootLineKeepSameInCurveList as rootLineKeepSameInCurveListRu,
  backgroundRepeatList as backgroundRepeatListRu,
  backgroundPositionList as backgroundPositionListRu,
  shapeList as shapeListRu,
  sidebarTriggerList as sidebarTriggerListRu,
  backgroundSizeList as backgroundSizeListRu,
  downTypeList as downTypeListRu,
  langList as langListRu
} from './ru'

const fontFamilyList = {
  zh: fontFamilyListZh,
  en: fontFamilyListEn,
  ru: fontFamilyListRu
}
const borderDasharrayList = {
  zh: borderDasharrayListZh,
  en: borderDasharrayListEn,
  ru: borderDasharrayListRu
}
const lineStyleList = {
  zh: lineStyleListZh,
  en: lineStyleListEn,
  ru: lineStyleListRu
}
const lineStyleMap = {
  zh: lineStyleMapZh,
  en: lineStyleMapZh
}
const rootLineKeepSameInCurveList = {
  zh: rootLineKeepSameInCurveListZh,
  en: rootLineKeepSameInCurveListEn,
  ru: rootLineKeepSameInCurveListRu
}
const backgroundRepeatList = {
  zh: backgroundRepeatListZh,
  en: backgroundRepeatListEn,
  ru: backgroundRepeatListRu
}
const backgroundPositionList = {
  zh: backgroundPositionListZh,
  en: backgroundPositionListEn,
  ru: backgroundPositionListRu
}
const backgroundSizeList = {
  zh: backgroundSizeListZh,
  en: backgroundSizeListEn,
  ru: backgroundSizeListRu
}
const shortcutKeyList = {
  zh: shortcutKeyListZh,
  en: shortcutKeyListEn
}
const shapeList = {
  zh: shapeListZh,
  en: shapeListEn,
  ru: shapeListRu
}

const shapeListMap = {
  zh: shapeListMapZh,
  en: shapeListMapZh
}

const sidebarTriggerList = {
  zh: sidebarTriggerListZh,
  en: sidebarTriggerListEn,
  ru: sidebarTriggerListRu
}

const downTypeList = {
  zh: downTypeListZh,
  en: downTypeListEn,
  ru: downTypeListRu
}

// merge language list to include ru without changing existing importers
const mergedLangList = Array.isArray(langList) ? Array.from(new Set([...(langList || []), ...(langListRu || [])])) : langList

export {
  fontSizeList,
  lineHeightList,
  borderWidthList,
  borderRadiusList,
  lineWidthList,
  store,
  colorList,
  mergedLangList as langList,
  fontFamilyList,
  borderDasharrayList,
  lineStyleList,
  lineStyleMap,
  rootLineKeepSameInCurveList,
  backgroundRepeatList,
  backgroundPositionList,
  backgroundSizeList,
  shortcutKeyList,
  shapeList,
  shapeListMap,
  sidebarTriggerList,
  downTypeList
}
