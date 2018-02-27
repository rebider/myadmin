<template>
  <el-main title="重置密码">
    <el-form ref="form" :model="form" :rules="rules" label-width="180px">
      <el-form-item label="旧密码" prop="oldPwd" >
        <el-input v-model="form.oldPwd" type="password"></el-input>
      </el-form-item>
      <el-form-item label="新密码" prop="newPwd">
        <el-input v-model="form.newPwd" type="password"></el-input>
      </el-form-item>
      <el-form-item label="重复新密码" prop="newPwd2">
        <el-input v-model="form.newPwd2" type="password"></el-input>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="onSubmit">修改</el-button>
      </el-form-item>
    </el-form>
  </el-main>
</template>
<script>
  import { changePasswd } from '@/api/user'
  import { Message } from 'element-ui'

  export default {
    data() {
      const validatePass = (rule, value, callback) => {
        if (value.length < 5) {
          callback(new Error('密码不能小于5位'))
        } else {
          callback()
        }
      }
      return {
        form: {
          oldPwd: '',
          newPwd: '',
          newPwd2: ''
        },
        rules: {
          oldPwd: [{ required: true, message: '请输入原密码', trigger: 'blur' }],
          newPwd: [{ required: true, trigger: 'blur', validator: validatePass }],
          newPwd2: [{ required: true, trigger: 'blur', validator: validatePass }]
        }
      }
    },
    methods: {
      onSubmit() {
        this.$refs.form.validate(valid => {
          if (valid) {
            if (this.form.newPwd !== this.form.newPwd2) {
              console.log('9999999999999999999')
              Message({
                message: '两次输入密码不一致!',
                type: 'error',
                duration: 2 * 1000
              })
              // this.$message({ type: 'error', message: '两次输入密码不一致' })
              return
            }
            console.log('oldPwd:', this.form.oldPwd)
            console.log('newPwd:', this.form.newPwd)
            console.log('newPwd2:', this.form.newPwd2)
            changePasswd(this.form).then(respone => {
              console.log(respone)
              Message({
                message: '修改密码成功!',
                type: 'info',
                duration: 3 * 1000
              })
            }).catch(error => {
              console.log(error)
            })
          } else {
            return false
          }
        })
      }
    }
  }
</script>
