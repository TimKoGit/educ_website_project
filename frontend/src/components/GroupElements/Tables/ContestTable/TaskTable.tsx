import * as React from 'react';
import axios from 'axios';
import './TaskTable.css';
import { useAppSelector, useAppDispatch } from '../../../../redux/hooks';
import { deleteTask } from '../../../../redux/slices/tasksSlice';


interface Task {
  id: string;
  name: string;
  url: string;
  status: string;
}


interface TaskTableProps {
  tasks: Task[];
}


const TaskTable: React.FC<TaskTableProps> = ({ tasks }) => {
  const [taskStatuses, setTaskStatuses] = React.useState<{ [key: string]: string }>(
    tasks.reduce((acc, task) => {
      acc[task.id] = task.status;
      return acc;
    }, {} as { [key: string]: string })
  );

  const fileInputRefs = React.useRef<{ [key: string]: HTMLInputElement | null }>({});

  const user = useAppSelector((state) => state.auth.user);
  const isStudent = user?.role === 'student';
  const isTeacher = user?.role === 'teacher';
  const dispatch = useAppDispatch();

  const handleUploadClick = (taskId: string) => {
    if (fileInputRefs.current[taskId]) {
      fileInputRefs.current[taskId]?.click();
    }
  };

  const handleFileChange = async (event: React.ChangeEvent<HTMLInputElement>, taskId: string) => {
    const file = event.target.files?.[0];
    if (file) {
      if (file.type !== 'image/png') {
        alert('Пожалуйста, выберите PNG файл.');
        return;
      }
      if (file.size > 2 * 1024 * 1024) { // 2MB limit
        alert('Размер файла должен быть меньше 2MB.');
        return;
      }

      const formData = new FormData();
      formData.append('file', file);

      try {
        await axios.post(`http://localhost:5000/tasks/${taskId}/upload`, formData, {
          headers: {
            'Content-Type': 'multipart/form-data',
          },
          withCredentials: true,
        });
        alert('Картинка отправлена успешно.');
        setTaskStatuses((prevStatuses) => ({
          ...prevStatuses,
          [taskId]: 'unchecked',
        }));
      } catch (error: any) {
        alert(`Не получилось отправить картинку: ${error.message}`);
      }
    }
  };

  const handleDelete = (taskId: string) => {
      if (window.confirm('Вы уверены, что хотите удалить эту задачу?')) {
        dispatch(deleteTask(taskId));
      }
    };

  const getStatusText = (status: string) => {
    console.log(status);
    switch (status) {
      case 'accepted':
        return <span style={{ color: 'green' }}>Принято</span>;
      case 'declined':
        return <span style={{ color: 'red' }}>Отклонено</span>;
      case 'unchecked':
        return <span style={{ color: 'black' }}>В проверке</span>;
      default:
        return <span style={{ color: 'black' }}>Не было посылок</span>;
    }
  };
    

  return (
    <div className='table-container contest-table-container'>
      <table>
        <thead>
          <tr>
            <th>Название задачи</th>
            <th>URL</th>
            <th>{isStudent && "Статус проверки"}</th>
            <th></th>
          </tr>
        </thead>
        <tbody>
          {tasks.map((task) => (
            <tr key={task.id} className='clickable-row'>
              <td>{task.name}</td>
              <td>
                <a href={task.url} target='_blank' rel='noopener noreferrer'>
                    {task.url}
                </a>
              </td>
              <td>{isStudent && getStatusText(taskStatuses[task.id])}</td>
              {isStudent && (
                <td>
                  <button className='upload-button' onClick={() => handleUploadClick(task.id)}>Загрузить изображение</button>
                  <input
                    type='file'
                    accept='image/png'
                    style={{ display: 'none' }}
                    ref={(el) => (fileInputRefs.current[task.id] = el)}
                    onChange={(e) => handleFileChange(e, task.id)}
                  />
                </td>
              )}
              {isTeacher && (
                <td>
                  <button className='delete-button' onClick={() => handleDelete(task.id)}>
                    Удалить задачу
                  </button>
                </td>
              )}
            </tr>
          ))}
        </tbody>
      </table>
    </div>
  );
};


export default TaskTable;
