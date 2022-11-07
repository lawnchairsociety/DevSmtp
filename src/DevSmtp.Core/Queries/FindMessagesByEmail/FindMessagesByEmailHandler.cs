using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Queries
{
    public sealed class FindMessagesByEmailHandler : IQueryHandler<FindMessagesByEmail, FindMessagesByEmailResult>
    {
        private readonly IDataStore _dataStore;

        public FindMessagesByEmailHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public async Task<FindMessagesByEmailResult> ExecuteAsync(FindMessagesByEmail query, CancellationToken cancellationToken = default)
        {
            cancellationToken.ThrowIfCancellationRequested();

            try
            {
                var messages = await this._dataStore.FindByEmailAsync(query.Email, cancellationToken);
                return new FindMessagesByEmailResult(messages);
            }
            catch (OperationCanceledException)
            {
                throw;
            }
            catch (Exception ex)
            {
                var message = $"Failed to find message by '{query.Email.Value}'.";
                var error = new FindMessagesByEmailException(message, ex);

                return new FindMessagesByEmailResult(error);
            }
        }
    }
}
