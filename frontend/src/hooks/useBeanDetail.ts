import { useState } from "react";
import { beansApi } from "@/api/beans";
import { showApiError } from "@/utils/errorHandler";
import type { CoffeeBean, UpdateBeanInput } from "@/types/api";

export interface UseBeanDetailReturn {
  bean: CoffeeBean | null;
  loading: boolean;
  editing: boolean;
  saving: boolean;
  loadBean: (id: string) => Promise<void>;
  updateBean: (id: string, input: UpdateBeanInput) => Promise<void>;
  deleteBean: (id: string) => Promise<boolean>;
  startEditing: () => void;
  cancelEditing: () => void;
  setBean: (bean: CoffeeBean) => void;
}

export function useBeanDetail(): UseBeanDetailReturn {
  const [bean, setBean] = useState<CoffeeBean | null>(null);
  const [loading, setLoading] = useState(true);
  const [editing, setEditing] = useState(false);
  const [saving, setSaving] = useState(false);

  const loadBean = async (id: string) => {
    try {
      const data = await beansApi.get(id);
      setBean(data);
    } catch (e) {
      showApiError(e, "データの取得に失敗しました");
    } finally {
      setLoading(false);
    }
  };

  const updateBean = async (id: string, input: UpdateBeanInput) => {
    setSaving(true);
    try {
      const updated = await beansApi.update(id, input);
      setBean(updated);
      setEditing(false);
    } catch (e) {
      showApiError(e, "更新に失敗しました");
    } finally {
      setSaving(false);
    }
  };

  const deleteBean = async (id: string): Promise<boolean> => {
    try {
      await beansApi.delete(id);
      return true;
    } catch (e) {
      showApiError(e, "削除に失敗しました");
      return false;
    }
  };

  const startEditing = () => setEditing(true);
  const cancelEditing = () => setEditing(false);

  return {
    bean,
    loading,
    editing,
    saving,
    loadBean,
    updateBean,
    deleteBean,
    startEditing,
    cancelEditing,
    setBean,
  };
}
