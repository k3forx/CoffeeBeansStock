type ValidationErrors = Partial<
  Record<"name" | "roastLevel" | "currentStock", string>
>;
type ValidationResult =
  | { valid: false; errors: ValidationErrors }
  | { valid: true; stock: number };

export function validateBeanForm(params: {
  name: string;
  roastLevel: string;
  currentStock: string;
}): ValidationResult {
  const errors: ValidationErrors = {};

  if (!params.name) {
    errors.name = "名前は必須です";
  }
  if (!params.roastLevel) {
    errors.roastLevel = "焙煎度を選択してください";
  }
  const stock = parseInt(params.currentStock, 10);
  if (isNaN(stock) || stock < 0) {
    errors.currentStock = "在庫数は0以上の数値を入力してください";
  }

  if (Object.keys(errors).length > 0) {
    return { valid: false, errors };
  }
  return { valid: true, stock };
}
