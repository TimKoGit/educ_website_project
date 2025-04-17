import * as React from 'react';
import { useState } from 'react';
import { useAppDispatch } from '../../../../redux/hooks';
import { addNewGroup } from '../../../../redux/slices/authSlice';
import './GroupForm.css';
import { useAppSelector } from '../../../../redux/hooks';
import { fetchGroupsByUserId } from '../../../../redux/slices/authSlice';


interface GroupFormProps {
  closeForm: () => void;
}


const GroupForm: React.FC<GroupFormProps> = ({ closeForm }) => {
  const dispatch = useAppDispatch();
  const [groupName, setGroupName] = useState('');
  const [groupCode, setGroupCode] = useState('');
  const user = useAppSelector((state) => state.auth.user);
  const [error, setError] = useState<string | null>(null);
  const { addNewGroupLoading, addNewGroupError } = useAppSelector((state) => state.auth);

  const handleFormSubmit = async (e: React.FormEvent) => {
    e.preventDefault();

    if (!groupName || !groupCode) {
      setError('Оба поля должны быть заполнены');
      return;
    }

    setError(null);

    try {
      await dispatch(addNewGroup({ name: groupName, code: groupCode })).unwrap();
      if (!addNewGroupError) {
      closeForm();
      setGroupName('');
      setGroupCode('');
      }
      debugger;
    } catch (error) {
    }

    try {
      await dispatch(fetchGroupsByUserId()).unwrap();
    } catch (error) {
      console.error('При получении групп появились ошибки:', error);
    }
  };

  const handleInvalid = (e: React.FormEvent<HTMLInputElement>) => {
    e.preventDefault();
    if (!groupName) {
      setError('Введите название группы');
    } else if (!groupCode) {
      setError('Введите код');
    }
  };

  return (
    <div className='group-form'>
      <h2>Новая группа!</h2>
      {error && <p style={{ color: 'red' }}>{error}</p>}
      {addNewGroupError && <p style={{ color: 'red' }}>{addNewGroupError}</p>}
      <form onSubmit={handleFormSubmit}>
        <div>
          <label>Название группы:</label>
          <input
            type='text'
            value={groupName}
            onChange={(e) => setGroupName(e.target.value)}
            required
            onInvalid={(e) => handleInvalid(e)}
          />
        </div>
        <div>
          <label>Код группы:</label>
          <input
            type='text'
            value={groupCode}
            onChange={(e) => setGroupCode(e.target.value)}
            required
            onInvalid={(e) => handleInvalid(e)}
          />
        </div>
        <button type='submit' disabled={addNewGroupLoading}>
          {addNewGroupLoading ? 'Создание...' : 'Создать'}
        </button>
      </form>
    </div>
  );
};


export default GroupForm;
