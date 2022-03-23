import React from 'react';
import { useForm } from 'react-hook-form';
import {
  FormControl,
  Input,
  TextField,
  InputLabel,
  FormHelperText,
  FormLabel,
  FormGroup,
  FormControlLabel,
  InputBase,
  Button,
  Checkbox,
  Switch,
  RadioGroup,
  Radio,
} from '@mui/material';

// 定义表单数据
interface IPerson {
  firstName: string;
  lastName: string;
}

export function Switches(): JSX.Element {
  const {
    register, // 注册字段
    handleSubmit, // 处理 submit 事件
    watch, // 实时监听数据的改变, 底层用了 useState(), 每次数据改变会刷新整个组件
    formState: { errors }, // 处理错误
    reset, // 重置表单
  } = useForm<IPerson>();

  function mySubmit(data: IPerson, event?: React.BaseSyntheticEvent): void {
    console.log(JSON.stringify(data, null, 2));
    // console.log(event);
  }

  return (
    <>
      <form onSubmit={handleSubmit(mySubmit)}>
        <FormControl component="fieldset">
          <FormLabel component="legend">formlabel</FormLabel>
          {/* register 到 RadioGroup 中 */}
          <RadioGroup row {...register('firstName')}>
            <FormControlLabel
              labelPlacement="top"
              label="form-control-label"
              control={<Radio />}
              value="top" // radio 的 value 放到 FormControlLabel 中
            />
            <FormControlLabel
              labelPlacement="bottom"
              label="form-control-label"
              control={<Radio />}
              value="bottom" // radio 的 value 放到 FormControlLabel 中
            />
          </RadioGroup>
        </FormControl>
        <Button type="submit">submit</Button>
      </form>
    </>
  );
}

// switch 是从 checkbox 衍生出来的, 标准 HTML input 中没有 switch 类型
export function MyTextField(): JSX.Element {
  const {
    register, // 注册字段
    handleSubmit, // 处理 submit 事件
    watch, // 实时监听数据的改变, 底层用了 useState(), 每次数据改变会刷新整个组件
    formState: { errors }, // 处理错误
    reset, // 重置表单
  } = useForm<IPerson>();

  function mySubmit(data: IPerson, event?: React.BaseSyntheticEvent): void {
    console.log(JSON.stringify(data, null, 2));
    // console.log(event);
  }

  return (
    <>
      <form onSubmit={handleSubmit(mySubmit)}>
        <FormControl component="fieldset">
          <FormLabel component="legend">formlabel</FormLabel>
          {/* FormGroup 可以控制组的元素排列 */}
          <FormGroup row>
            <FormControlLabel
              labelPlacement="top"
              label="form-control-label"
              control={<Checkbox value="top" {...register('firstName')} />}
            />
            <FormControlLabel labelPlacement="bottom" label="form-control-label" control={<Checkbox value="bot" />} />
          </FormGroup>
        </FormControl>
        <Button type="submit">submit</Button>
      </form>
    </>
  );
}

export function TestFormControl(): JSX.Element {
  return (
    <div style={{ margin: '20px' }}>
      <FormControl>
        <FormLabel>FormLabel</FormLabel>
        <FormGroup row>
          <FormControlLabel
            label="top" // label 的名字
            labelPlacement="top" // label 的位置
            control={<Checkbox value="top" />} // 需要 label 的组件
          />
          <FormControlLabel
            label="start" // label 的名字
            labelPlacement="start" // label 的位置
            control={<Checkbox value="start" />} // 需要 label 的组件
          />
          <FormControlLabel
            label="bottom" // label 的名字
            labelPlacement="bottom" // label 的位置
            control={<Checkbox value="bottom" />} // 需要 label 的组件
          />
          <FormControlLabel
            label="end" // label 的名字
            labelPlacement="end" // label 的位置
            control={<input value="end" />} // 需要 label 的组件
          />
        </FormGroup>
      </FormControl>
    </div>
  );
}

export function CheckBoxSelect(): JSX.Element {
  const {
    register, // 注册字段
    handleSubmit, // 处理 submit 事件
    watch, // 实时监听数据的改变, 底层用了 useState(), 每次数据改变会刷新整个组件
    formState: { errors }, // 处理错误
    reset, // 重置表单
  } = useForm<IPerson>();

  function mySubmit(data: IPerson, event?: React.BaseSyntheticEvent): void {
    console.log(JSON.stringify(data, null, 2));
    // console.log(event);
  }

  return (
    <>
      <form onSubmit={handleSubmit(mySubmit)}>
        <FormControl component="fieldset">
          <FormLabel component="legend">formlabel</FormLabel>

          {/* ⭐️ register 到 RadioGroup 中 */}
          <FormGroup row>
            <FormControlLabel
              labelPlacement="top"
              label="A"
              control={<Checkbox value="A" {...register('firstName')} />}
              // ⭐️ radio 的 value 放到 FormControlLabel 中
            />
            <FormControlLabel
              labelPlacement="bottom"
              label="B"
              control={<Checkbox value="B" {...register('firstName')} />}
              // ⭐️ radio 的 value 放到 FormControlLabel 中
            />
          </FormGroup>
        </FormControl>
        <Button type="submit">submit</Button>
      </form>
    </>
  );
}
