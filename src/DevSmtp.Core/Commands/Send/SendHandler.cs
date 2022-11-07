using DevSmtp.Core.Stores;

namespace DevSmtp.Core.Commands
{
    public sealed class SendHandler : ICommandHandler<Send, SendResult>
    {
        private readonly IDataStore _dataStore;

        public SendHandler(IDataStore dataStore)
        {
            this._dataStore = dataStore ?? throw new ArgumentNullException(nameof(dataStore));
        }

        public Task<SendResult> ExecuteAsync(Send command, CancellationToken cancellationToken = default)
        {
            throw new NotImplementedException();
        }
    }
}
