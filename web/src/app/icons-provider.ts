import { EnvironmentProviders, importProvidersFrom } from '@angular/core';
import {
  MenuFoldOutline,
  MenuUnfoldOutline,
  PlusCircleOutline,
  QuestionCircleOutline,
  ExclamationCircleOutline,
  TeamOutline,
  UserOutline,
  UploadOutline,
  FileAddOutline,
  LockOutline,
  ExpandAltOutline,
  EditOutline,
  DeleteOutline,
  CopyOutline,
} from '@ant-design/icons-angular/icons';
import { NzIconModule } from 'ng-zorro-antd/icon';

const icons = [
  MenuFoldOutline,
  MenuUnfoldOutline,
  PlusCircleOutline,
  QuestionCircleOutline,
  ExclamationCircleOutline,
  TeamOutline,
  UserOutline,
  UploadOutline,
  FileAddOutline,
  LockOutline,
  ExpandAltOutline,
  EditOutline,
  DeleteOutline,
  CopyOutline,
];

export function provideNzIcons(): EnvironmentProviders {
  return importProvidersFrom(NzIconModule.forRoot(icons));
}
