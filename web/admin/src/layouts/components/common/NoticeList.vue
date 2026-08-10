<!-- Author: Charlie -->

<script setup lang="ts">
export interface BannerItem {
  avatar?: string | null
  id: string
  type?: number
  title: string
  icon: string
  tagTitle?: string
  /** 字典色，优先于 tagType */
  tagColor?: { color: string; textColor?: string }
  tagType?: 'default' | 'error' | 'primary' | 'info' | 'success' | 'warning'
  description?: string
  date: string
  isRead?: boolean
}

const avatarImgProps = { referrerPolicy: 'no-referrer' } as any

defineProps<{
  list?: BannerItem[]
  loading?: boolean
  hasMore?: boolean
}>()

const emit = defineEmits<{
  open: [id: string]
  loadMore: []
}>()
</script>

<template>
  <NScrollbar style="height: 360px">
    <NEmpty
      v-if="!loading && !list?.length"
      description="暂无消息"
      style="padding: 64px 0"
    />
    <NSpace
      v-else-if="loading && !list?.length"
      justify="center"
      style="padding: 120px 0"
    >
      <NSpin size="small" />
    </NSpace>
    <NList
      v-else
      hoverable
      clickable
    >
      <NListItem
        v-for="item in list"
        :key="item.id"
        @click="emit('open', item.id)"
      >
        <NThing>
          <template #avatar>
            <NBadge
              :dot="!item.isRead"
              :processing="!item.isRead"
              type="info"
            >
              <NAvatar
                v-if="item.avatar"
                round
                :size="32"
                :src="item.avatar || undefined"
                :img-props="avatarImgProps"
              />
              <NAvatar
                v-else
                round
                :size="32"
              >
                <NovaIcon
                  :icon="item.icon"
                  :size="16"
                  :style="{
                    color: item.isRead ? 'var(--text-color-3)' : 'var(--primary-color)',
                  }"
                />
              </NAvatar>
            </NBadge>
          </template>
          <template #header>
            <NEllipsis style="max-width: 220px">
              <NText
                :depth="item.isRead ? 3 : 1"
                :strong="!item.isRead"
              >
                {{ item.title }}
              </NText>
            </NEllipsis>
          </template>
          <template #header-extra>
            <NTag
              v-if="item.tagTitle"
              size="tiny"
              :bordered="false"
              :color="item.tagColor"
              :type="item.tagColor ? undefined : item.tagType || 'default'"
            >
              {{ item.tagTitle }}
            </NTag>
          </template>
          <template #description>
            <div>
              <NEllipsis
                v-if="item.description"
                :line-clamp="1"
                :tooltip="false"
              >
                <NText
                  depth="3"
                  style="font-size: 12px"
                >
                  {{ item.description }}
                </NText>
              </NEllipsis>
              <NText
                depth="3"
                style="font-size: 11px; display: block"
              >
                {{ item.date }}
              </NText>
            </div>
          </template>
        </NThing>
      </NListItem>
      <NSpace
        v-if="hasMore"
        justify="center"
        style="padding: 8px 0 12px"
      >
        <NButton
          text
          size="small"
          :loading="loading"
          @click.stop="emit('loadMore')"
        >
          加载更多
        </NButton>
      </NSpace>
    </NList>
  </NScrollbar>
</template>
