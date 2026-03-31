import { useState } from "react";
import type {
  CoffeeBean,
  CreateBeanInput,
  UpdateBeanInput,
  RoastLevel,
  RoastDetail,
} from "@/types/api";
import { validateBeanForm } from "@/utils/validation";

type BeanFormFields = {
  name: string;
  origin: string;
  roastLevel: RoastLevel | "";
  roastDetail: RoastDetail | "";
  currentStock: string;
  notes: string;
};

type BeanFormField = keyof BeanFormFields;

export interface UseBeanFormReturn {
  fields: BeanFormFields;
  errors: Record<string, string>;
  setField: <K extends BeanFormField>(field: K, value: BeanFormFields[K]) => void;
  handleRoastLevelSelect: (level: RoastLevel) => void;
  validate: () => { valid: true; stock: number } | { valid: false; errors: Record<string, string> };
  resetToBean: (bean: CoffeeBean) => void;
  toCreateInput: (stock: number) => CreateBeanInput;
  toUpdateInput: (stock: number) => UpdateBeanInput;
}

export function useBeanForm(): UseBeanFormReturn {
  const [name, setName] = useState("");
  const [origin, setOrigin] = useState("");
  const [roastLevel, setRoastLevel] = useState<RoastLevel | "">("");
  const [roastDetail, setRoastDetail] = useState<RoastDetail | "">("");
  const [currentStock, setCurrentStock] = useState("");
  const [notes, setNotes] = useState("");
  const [errors, setErrors] = useState<Record<string, string>>({});

  const fieldSetters: Record<BeanFormField, (v: string) => void> = {
    name: setName,
    origin: setOrigin,
    roastLevel: setRoastLevel as (v: string) => void,
    roastDetail: setRoastDetail as (v: string) => void,
    currentStock: setCurrentStock,
    notes: setNotes,
  };

  const setField = <K extends BeanFormField>(field: K, value: BeanFormFields[K]) => {
    fieldSetters[field](value as string);
    setErrors((prev) => {
      const { [field]: _, ...rest } = prev;
      return rest;
    });
  };

  const handleRoastLevelSelect = (level: RoastLevel) => {
    setRoastLevel(level);
    setRoastDetail("");
  };

  const validate = () => {
    const result = validateBeanForm({ name, roastLevel, currentStock });
    if (!result.valid) {
      setErrors(result.errors as Record<string, string>);
      return { valid: false as const, errors: result.errors as Record<string, string> };
    }
    setErrors({});
    return { valid: true as const, stock: result.stock };
  };

  const resetToBean = (bean: CoffeeBean) => {
    setName(bean.name);
    setOrigin(bean.origin ?? "");
    setRoastLevel(bean.roast_level);
    setRoastDetail(bean.roast_detail ?? "");
    setCurrentStock(String(bean.current_stock));
    setNotes(bean.notes ?? "");
    setErrors({});
  };

  const toCreateInput = (stock: number): CreateBeanInput => ({
    name,
    origin: origin || undefined,
    roast_level: roastLevel as RoastLevel,
    roast_detail: roastDetail || undefined,
    current_stock: stock,
    notes: notes || undefined,
  });

  const toUpdateInput = (stock: number): UpdateBeanInput => ({
    name,
    origin: origin || undefined,
    roast_level: roastLevel as RoastLevel,
    roast_detail: roastDetail || undefined,
    current_stock: stock,
    notes: notes || undefined,
  });

  return {
    fields: { name, origin, roastLevel, roastDetail, currentStock, notes },
    errors,
    setField,
    handleRoastLevelSelect,
    validate,
    resetToBean,
    toCreateInput,
    toUpdateInput,
  };
}
