import { Alert } from "react-native";

export function validateBeanForm(params: {
  name: string;
  roastLevel: string;
  currentStock: string;
}): { valid: false } | { valid: true; stock: number } {
  if (!params.name) {
    Alert.alert("エラー", "名前は必須です");
    return { valid: false };
  }
  if (!params.roastLevel) {
    Alert.alert("エラー", "焙煎度を選択してください");
    return { valid: false };
  }
  const stock = parseInt(params.currentStock, 10);
  if (isNaN(stock) || stock < 0) {
    Alert.alert("エラー", "在庫数は0以上の数値を入力してください");
    return { valid: false };
  }
  return { valid: true, stock };
}
