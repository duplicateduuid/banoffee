<script lang="ts">
	import classNames from 'classnames';

	type Props = {
		redirect: string;
		name: string;
		description?: string | null;
		author?: string | null;
		className?: string;
	};
	const { redirect, name, description, author, className }: Props = $props();
</script>

<a
	href={redirect}
	class={classNames('flex gap-4 h-40 shadow px-8 py-6 rounded-lg bg-white', className)}
>
	<span class="bg-primary-400 h-full aspect-video rounded-lg" />
	<div class="flex flex-col justify-between h-full">
		<div class="flex flex-col">
			<p class="font-primary text-2xl tracking-wider font-semibold">
				{name}
			</p>
			{#if description}
				<p class="text-sm line-clamp-2">{description}</p>
			{/if}
		</div>
		{#if author}
			<p>{author}</p>
		{/if}
	</div>
</a>
